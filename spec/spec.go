package spec

import (
  "bytes"
  "encoding/json"
	"errors"
  "fmt"
	"github.com/paddyforan/jarvis/parse"
	"io"
	"strings"
)

var UnsupportedOutputFormatError = errors.New("Unsupported output format.")

func Generate(outputFormat string, output io.WriteCloser, resources []*parse.Resource) error {
	defer output.Close()
	for _, resource := range resources {
		if resource == nil {
			continue
		}
		endpoints, err := BuildEndpoints(*resource)
		if err != nil {
			return err
		}
		if len(endpoints) < 1 {
			continue
		}
		err = writeResourceHeader(output, outputFormat, resource, resource.ID)
		if err != nil {
			return err
		}
		err = writeProperties(output, outputFormat, resource)
		if err != nil {
			return err
		}
		for _, endpoint := range endpoints {
			err = writeEndpoint(output, outputFormat, endpoint)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeResourceHeader(output io.Writer, outputFormat string, resource *parse.Resource, id string) error {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "markdown":
		_, err := fmt.Fprintf(output, "\n# %s (%s)\n%s", resource.Name, id, resource.Description)
		return err
	default:
		return UnsupportedOutputFormatError
	}
}

func writeProperties(output io.Writer, outputFormat string, resource *parse.Resource) error {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "markdown":
    if len(resource.Properties) > 0 {
      _, err := fmt.Fprint(output, "\n")
      if err != nil {
        return err
      }
    }
		for _, property := range resource.Properties {
			_, err := fmt.Fprintf(output, "\n * **%s** *(%s)*: %s", property.ID, property.Type, property.Description)
			if err != nil {
				return err
			}
			if len(property.Values) > 0 {
				_, err = fmt.Fprint(output, "\n\t * **Possible Values**:")
				for _, value := range property.Values {
					_, err = fmt.Fprintf(output, "\n\t\t * %v", value)
				}
			}
			if property.Default != nil {
				_, err = fmt.Fprintf(output, "\n\t * **Default Value**: %v", property.Default)
			}
			if property.Maximum != 0 {
				_, err = fmt.Fprintf(output, "\n\t * **Maximum Value**: %v", property.Maximum)
			}
			if property.Minimum != 0 {
				_, err = fmt.Fprintf(output, "\n\t * **Minimum Value**: %v", property.Minimum)
			}
		}
	default:
		return UnsupportedOutputFormatError
	}
	return nil
}

func writeEndpoint(output io.Writer, outputFormat string, endpoint Endpoint) error {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "markdown":
		querystring := ""
		for _, param := range endpoint.Params {
			if param.Default != nil {
				continue
			}
			if querystring != "" {
				querystring += "&"
			}
			querystring += param.ID + "={" + param.Type + "}"
			if param.Repeated {
				querystring += "&" + param.ID + "={" + param.Type + "}"
				querystring += "&" + param.ID + "={" + param.Type + "}"
			}
		}
		if querystring != "" {
			querystring = "?" + querystring
		}
		_, err := fmt.Fprintf(output, "\n\n## %s\n\n### Request\n\n%s /%s%s", endpoint.Name, endpoint.Verb, endpoint.Path, querystring)
		if err != nil {
			return err
		}
		if len(endpoint.SampleRequest) > 0 {
			_, err := fmt.Fprint(output, "\n\n\t")
			if err != nil {
				return err
			}
      buf := bytes.NewBuffer([]byte{})
			err = json.Indent(buf, endpoint.SampleRequest, "\t", "  ")
			if err != nil {
				return err
			}
      _, err = buf.WriteTo(output)
      if err != nil {
        return err
      }
		}
	default:
		return UnsupportedOutputFormatError
	}
	return nil
}
