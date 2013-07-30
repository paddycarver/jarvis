# API Principles

The following principles are built into the `apidef` command and are considered requirements for the API. While subverting them may be possible, they are enforced as the default.

## One representation to rule them all

Each resource gets one and only one representation. This consistent representation must be used every time the resource is returned or accepted by the API.

Returning IDs instead of embedding full resources is acceptable, but the field name must make it clear that the _ID_ is what is being represented, not the resource. For example, relations should use a `_id` suffix to denote a field that represents the ID.

## Actionable errors

Each error in the 4XX set of codes must have one and only one action necessary to resolve it. Errors should precisely point to the field that caused the error and be utilised as a way to explain what went wrong. For example, &quot;unacceptable ID&quot; is bad; &quot;ID too long&quot; is good.

Errors must also be machine readable; no information that must be used to identify the error can be returned in a string unless the string is an ID.

## HTTP Rules

Gets are idempotent. Posts create new resources. Puts replace entire resources. Patches replace specific parts of a resource. Deletes remove a resource.
