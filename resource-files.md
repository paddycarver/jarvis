# Resource Files

Resource files are machine-readable files that define the resources made available by the APIs being generated. They are JSON documents that describe the resources, along with their interactions and properties.

## Sample

	:::json
	{
	  "_version": "0.1.0",
	  "id": "message",
	  "name": "Message",
	  "description": "A small chunk of data that is meant to be processed.",
	  "parent": "queue",
	  "url_slug": "messages",
	  "properties": [
	    {
	      "id": "id",
	      "type": "string",
	      "description": "A unique identifier for the message.",
	      "required": false,
	      "format": "{regular expression}"
	    },
	    {
	      "id": "body",
	      "type": "binary",
	      "description": "The data that is meant to be processed.",
	      "required": true,
	      "maximum": 65536
	    }
	  ],
	  "interactions": [
	    {
	      "id": "push",
	      "verb": "create",
	      "description": "Push a message onto the end of the queue.",
	      "omitted_input_fields": ["id"]
	    }
	  ]
	}

## Properties

The properties of the resource object determine the resource's behaviour and properties.

<table>
<tr><th>Field</th><th>Required</th><th>Description</th></tr>
<tr><td>_version</td><td>Yes</td><td>The version of the resource file syntax to use when reading the syntax file.</td></tr>
<tr><td>id</td><td>Yes</td><td>An API-unique ID for the resource.</td></tr>
<tr><td>name</td><td>Yes</td><td>A human-friendly name for the resource.</td></tr>
<tr><td>description</td><td>Yes</td><td>A human-friendly description of the resource.</td></tr>
<tr><td>parent</td><td>No</td><td>The ID of the resource this resource is a child of, if this resource has a parent.</td></tr>
<tr><td>url_slug</td><td>Yes</td><td>The slug that will be used when constructing URLs for this resource.</td></tr>
<tr><td>properties</td><td>Yes</td><td>Property objects describing the properties of the resource.</td></tr>
<tr><td>interactions</td><td>No</td><td>Interaction objects describing the possible actions that can be performed against the resource.</td></tr>
</table>

Property objects have their own properties, describing the constraints of the property:

<table>
<tr><th>Field</th><th>Required</th><th>Description</th></tr>
<tr><td>id</td><td>Yes</td><td>A resource-unique ID for the property.</td></tr>
<tr><td>type</td><td>Yes</td><td>The type of value expected by the property. Should be one of the following: string, bytes, time, datetime, int, float, boolean, array, object, pointer</td></tr>
<tr><td>description</td><td>Yes</td><td>A human-friendly description of the property.</td></tr>
<tr><td>required</td><td>Yes</td><td>If true, requests omitting the property will be considered invalid (unless specifically overridden in the interaction) and responses will return the property even if it is set to an empty value.</td></tr>
<tr><td>format</td><td>No</td><td>A regular expression that the value of the property must match. Requests with properties not matching this format will be considered invalid unless specifically overridden in the interaction.</td></tr>
<tr><td>maximum</td><td>No</td><td>A maximum value, as an int, for the value of the property. For strings, bytes, and arrays, the length ix compared to the maximum value. For times, datetimes, ints, and floats, the value is compared to the maximum value. Objects and pointers cannot have maximum values.</td></tr>
<tr><td>minimum</td><td>No</td><td>A minimum value, as an int, for the value of the property. For strings, bytes, and arrays, the length ix compared to the minimum value. For times, datetimes, ints, and floats, the value is compared to the minimum value. Objects and pointers cannot have minimum values.</td></tr>
<tr><td>value_type</td><td>No</td><td>For pointers, the type of the value the pointer is pointing to. Requests pointing to other types will be considered invalid unless specifically overridden in the interaction.</td></tr>
</table>

Interaction objects have their own properties, describing the constraints and requirements of the interaction:

<table>
<tr><th>Field</th><th>Required</th><th>Description</th></tr>
<tr><td>id</td><td>Yes</td><td>A resource-unique ID for the interaction.</td></tr>
<tr><td>verb</td><td>Yes</td><td>A description of what the interaction does to the resource. Accepted values are: create, read, update, replace, destroy</td></tr>
<tr><td>description</td><td>Yes</td><td>A human-friendly description of the interaction.</td></tr>
<tr><td>omitted_input_fields</td><td>No</td><td>A list of field IDs that will be ignored when a resource is used as input (as in a request).</td></tr>
<tr><td>rejected_input_fields</td><td>No</td><td>A list of field IDs that will, if included, cause input to be considered invalid.</td></tr>
<tr><td>required_input_fields</td><td>No</td><td>A list of field IDs that will, if not included, cause input to be considered invalid. These required fields are in addition to the fields whose required attribute is true.</td></tr>
<tr><td>omitted_output_fields</td><td>No</td><td>A list of field IDs that will be emptied before returning a resource as output (as in a response).</td></tr>
<tr><td>rejected_output_fields</td><td>No</td><td>A list of field IDs that will cause a server error if they are populated as part of output.</td></tr>
<tr><td>required_output_fields</td><td>No</td><td>A list of field IDs that will cause a server error if they cannot be populated as part of output.</td></tr>
</table>
