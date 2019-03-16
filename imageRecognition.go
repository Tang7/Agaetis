package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

const (
	graphFile  = "./model/inception5h/tensorflow_inception_graph.pb"
	labelsFile = "./model/inception5h/imagenet_comp_graph_label_strings.txt"
)

// ErrInvalidImageURL means user input an invalid image URL
var ErrInvalidImageURL = errors.New("Unable to get image from URL")

// ErrLoadModelFailed means app failed to load pre-trained model
var ErrLoadModelFailed = errors.New("Load model failed")

// ErrImageNormalizeFailed means app failed to normalized the image
var ErrImageNormalizeFailed = errors.New("Image normalize failed")

// ErrSessionInitFailed means unable to initilze the tensorflow session
var ErrSessionInitFailed = errors.New("Session init failed")

// ErrSessionRunFailed means error happens during session running
var ErrSessionRunFailed = errors.New("Session run failed")

// SuccessImageRecognition is success image recognition text
var SuccessImageRecognition = "Image Recognition Done"

// ImageDetails test use
type ImageDetails struct {
	URL           string
	PredLabels    Labels
	Success       bool
	StatusMessage string
}

// Label store preditive result as label and probability
type Label struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

// Labels predition for input image
type Labels []Label

// Support func to sort defined struct Labels
func (l Labels) Len() int      { return len(l) }
func (l Labels) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Less compare operation to sort Labels based on Probability
func (l Labels) Less(i, j int) bool { return l[i].Probability > l[j].Probability }

// UploadImage test use
func UploadImage(w http.ResponseWriter, r *http.Request) {
	ImageRecognitionPage := template.Must(template.ParseFiles("html/ImageRecognition.html"))
	if r.Method != "POST" {
		ImageRecognitionPage.Execute(w, nil)
		return
	}

	details := ImageDetails{URL: r.FormValue("url")}

	var err error
	details.PredLabels, err = makePrediction(details.URL)
	if err != nil {
		details.StatusMessage = err.Error()
		ImageRecognitionPage.Execute(w, details)
		return
	}

	if len(details.PredLabels) > 0 {
		details.StatusMessage = SuccessImageRecognition
		details.Success = true
	}

	ImageRecognitionPage.Execute(w, details)
}

func makePrediction(url string) (Labels, error) {

	response, e := http.Get(url)
	if e != nil {
		return nil, ErrInvalidImageURL
	}
	defer response.Body.Close()

	// load model
	graph, labels, err := loadModel()
	if err != nil {
		return nil, ErrLoadModelFailed
	}

	tensor, err := normalizeImage(response.Body)
	if err != nil {
		return nil, ErrImageNormalizeFailed
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, ErrSessionInitFailed
	}

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("output").Output(0),
		},
		nil)

	if err != nil {
		return nil, ErrSessionRunFailed
	}

	predLabels := getTopFiveLabels(labels, output[0].Value().([][]float32)[0])
	for _, l := range predLabels {
		fmt.Printf("label: %s, probability: %.2f%%\n", l.Label, l.Probability*100)
	}
	return predLabels, nil
}

func loadModel() (*tf.Graph, []string, error) {
	model, err := ioutil.ReadFile(graphFile)
	if err != nil {
		return nil, nil, err
	}
	// initial graph, and import pre-trained model into it
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return nil, nil, err
	}
	// load and store labels
	openLabelsFile, err := os.Open(labelsFile)
	if err != nil {
		return nil, nil, err
	}
	defer openLabelsFile.Close()
	scanner := bufio.NewScanner(openLabelsFile)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	return graph, labels, scanner.Err()
}

func normalizeImage(img io.ReadCloser) (*tf.Tensor, error) {
	var buf bytes.Buffer
	// TODO check image size before copy
	_, err := io.Copy(&buf, img)
	if err != nil {
		return nil, err
	}
	// convert go value to Tensor
	tensor, err := tf.NewTensor(buf.String())
	if err != nil {
		return nil, err
	}
	// normalized data to be 224x224 before input into model
	graph, input, output, err := getNormalizedGraph()
	if err != nil {
		return nil, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	normalizedImg, err := session.Run(
		map[tf.Output]*tf.Tensor{
			input: tensor,
		},
		[]tf.Output{
			output,
		},
		nil)

	if err != nil {
		return nil, err
	}

	return normalizedImg[0], nil
}

// getNormalizedGraph decode, rezise and normalize input image
func getNormalizedGraph() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	// decode RGB image
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	// Sub: returns x - y element-wise
	output = op.Sub(s,
		// make it 224x224: inception specific
		op.ResizeBilinear(s,
			// inserts a dimension of 1 into a tensor's shape.
			op.ExpandDims(s,
				// cast image to float type
				op.Cast(s, decode, tf.Float),
				op.Const(s.SubScope("make_batch"), int32(0))),
			op.Const(s.SubScope("size"), []int32{224, 224})),
		// mean = 117: inception specific
		op.Const(s.SubScope("mean"), float32(117)))
	graph, err = s.Finalize()

	return graph, input, output, err
}

func getTopFiveLabels(labels []string, probabilities []float32) []Label {
	var resultLabels []Label
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, Label{Label: labels[i], Probability: p})
	}

	sort.Sort(Labels(resultLabels))
	return resultLabels[:5]
}
