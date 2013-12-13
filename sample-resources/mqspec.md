
# Message (message)
A string of data that is meant to be processed.

 * **id** *(string)*: A unique, API-generated identifier for this resource.
 * **body** *(string)*: The data that is meant to be processed.
 * **timeout** *(duration)*: The amount of time a single reservation lasts for this message by default. This can be overwritten on a per-request basis.
	 * **Default Value**: 60
	 * **Maximum Value**: 86400
	 * **Minimum Value**: 30
 * **delay** *(duration)*: The number of seconds to delay putting a message on the queue. The message will not be available until this time has elapsed.
	 * **Default Value**: 0
	 * **Maximum Value**: 604800
 * **expires_in** *(duration)*: The number of seconds to keep a message on the queue before it is automatically deleted.
	 * **Default Value**: 604800
	 * **Maximum Value**: 2592000

## Delete a Message

### Request

DELETE /projects/{id}/queues/{name}/messages/{id}

## Delete Messages

### Request

DELETE /projects/{id}/queues/{name}/messages?id={string}&id={string}&id={string}

## Peek at Messages

### Request

GET /projects/{id}/queues/{name}/messages

## Push Messages

### Request

POST /projects/{id}/queues/{name}/messages

	{
	  "messages": [
	    {
	      "body": "",
	      "delay": 0,
	      "expires_in": 604800,
	      "timeout": 60
	    },
	    {
	      "body": "35RNWJNAd+jNFHT5xnM=",
	      "delay": 0
	    },
	    {
	      "body": "BO2pSdGakn2CoMHrTm8clTkThE5Z",
	      "delay": 0,
	      "expires_in": 604800
	    }
	  ]
	}
# Queue (queue)
An ordered set of messages that messages can be added to or pulled off of.

 * **id** *(string)*: A unique, API-generated identifier for this resource.
 * **name** *(string)*: A unique, human-readable identifier for the queue.
 * **push_type** *(string)*: How messages will be removed from the queue. Pull queues store messages until they're asked for by a client. Multicast queues send HTTP callbacks to each subscriber to deliver messages. Unicast queues send a single HTTP callback to a subscriber chosen at random.
	 * **Possible Values**:
		 * pull
		 * multicast
		 * unicast
	 * **Default Value**: pull
 * **retries** *(int)*: The number of times an HTTP callback should be retried for push queues.
	 * **Default Value**: 3
	 * **Maximum Value**: 100
 * **retries_delay** *(duration)*: The number of seconds each retry of the HTTP callback should be delayed for push queues.
	 * **Default Value**: 60
	 * **Maximum Value**: 86400
	 * **Minimum Value**: 30

## List Queues

### Request

GET /projects/{id}/queues

## Delete a Queue

### Request

DELETE /projects/{id}/queues/{name}

## Get Queue Info

### Request

GET /projects/{id}/queues/{name}

## Create Queue

### Request

POST /projects/{id}/queues

	{
	  "queue": {
	    "name": "+dIYOq2HrlcYbw==",
	    "push_type": "pull",
	    "retries": 3
	  }
	}

## Update Queue Info

### Request

PUT /projects/{id}/queues/{name}

	{
	  "queue": {
	    "name": "Cqi2go/g1L1q4Q==",
	    "push_type": "pull",
	    "retries": 3,
	    "retries_delay": 60
	  }
	}
# Reservation (reservation)
A lock on a message that prevents another client from retrieving the message for a short duration.

 * **id** *(string)*: A unique, API-generated identifier for this resource.
 * **message_id** *(string)*: A unique, API-generated identifier that points to the message this reservation is locking.
 * **timeout** *(duration)*: The amount of time this reservation lasts before the lock on the message will be released.
	 * **Default Value**: 60
	 * **Maximum Value**: 86400
	 * **Minimum Value**: 30

## Release a Message

### Request

DELETE /projects/{id}/queues/{name}/messages/reservations/{id}

## Get a Message

### Request

POST /projects/{id}/queues/{name}/messages/reservations

	{
	  "reservation": {
	    "timeout": 60
	  }
	}

## Touch a Message

### Request

PUT /projects/{id}/queues/{name}/messages/reservations/{id}

## Get Reservation Information

### Request

GET /projects/{id}/queues/{name}/messages/reservations/{id}
# Subscriber (subscriber)
A client interested in receiving push messages from this queue.

 * **url** *(string)*: The URL endpoint push messages should be POSTed to.

## Add Subscribers

### Request

POST /projects/{id}/queues/{name}/subscribers

	{
	  "subscribers": [
	    {
	      "url": "Z5ZWNTifOt2QItky/KffazJihZNfxlY="
	    },
	    {
	      "url": "rVPDSNPOGmwteyqPhJ8="
	    },
	    {
	      "url": "Gd0="
	    }
	  ]
	}

## Remove Subscribers

### Request

DELETE /projects/{id}/queues/{name}/subscribers

## List Subscribers

### Request

GET /projects/{id}/queues/{name}/subscribers
# Subscription (subscription)
The status information about a subscriber for each message.

 * **id** *(string)*: A unique, API-generated identifier for this resource.
 * **retries_delay** *(duration)*: When a push fails, this duration specifies the delay before the push will be retried.
 * **retries_remaining** *(int)*: The number of times a push message will be retried before it is considered failed permanently.
 * **status_code** *(int)*: The HTTP status code returned by the subscriber on the last push.
 * **status** *(string)*: A string describing the status of the subscription.
 * **url** *(string)*: The URL of the subscriber this subscription belongs to.

## List Subscriptions

### Request

GET /projects/{id}/queues/{name}/messages/{id}/subscriptions

## Acknowledge a Push Message

### Request

DELETE /projects/{id}/queues/{name}/messages/{id}/subscriptions/{id}
# Webhook (webhook)
An input that will create a message out of whatever request it receives.

## Push Message From Webhook

### Request

POST /projects/{id}/queues/{name}/messages/webhook?oauth={string}