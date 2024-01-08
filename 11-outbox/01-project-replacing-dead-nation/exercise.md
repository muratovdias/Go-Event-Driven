# Getting Rid of Dead Nation


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

Our company is getting bigger and bigger.
We are selling more and more tickets, and we are starting to hit the wall of the current setup.

Currently, we are using a company called "Dead Nation" to handle all logic related to ticket reservation, keeping track of the seat limits, etc.
For example, the well known `POST /tickets-status` endpoint is called by Dead Nation when a ticket is confirmed.

Dead Nation starts to become a problem for multiple reasons.
First of all, they don't support reserving individual seats, so we can't support big events that happen in stadiums.

Sometimes their system also does not work properly — for example, they might allow booking a ticket for a show that is sold out.
It's problematic to handle some exclusive events with big demand: Imagine selling 500 tickets for show limited to 50 people.

These issues are making our customers angry and generating a ton of work for our ops team.

Sometimes matters are even worse: Their system is down for long periods of time, especially if we have a big event.

There is also one more problem: They are not cheap.
Dead Nation is taking a significant cut from our margin.

Time to think about how we can migrate away from Dead Nation.
It would require a big engineering effort to do this, but we have a good (financial) reason: All of the previously described problems are costing us a lot of money.


Big migrations are always risky.
We should make this migration not a big bang but rather take it step by step.

</span>
	</div>
	</div>

## The Strangler Pattern

To migrate from Dead Nation, we will use a form of the [Strangler Pattern](https://www.martinfowler.com/bliki/StranglerFigApplication.html).

In the Strangler Pattern, we create a "Strangler" part of our application that is a facade in our system and proxies all incoming requests.
Later, the "Strangler" part delegates all requests both to new and old parts of our application.

Thanks to that, the legacy system works as before, but in the meantime, we can migrate all functionality to the new system.
At this stage, when we think that our new application is ready, we can do cross-checks between systems to ensure that none of the data was lost and everything works as expected.

This is how this usually works at a high level:

```plantuml
@startuml

skinparam monochrome true

!define ICONURL https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/v2.5.0/dist

!includeurl ICONURL/C4.puml

title Strangler Pattern

actor "User" as user

node "Legacy System" as legacy {
  component "Legacy Application" as legacyApp
}

node "Strangler Application" as strangler {
  database "Strangler Database" as stranglerDB
  component "Strangler Module" as stranglerModule
}

node "New Application" as newApp {
  database "New Database" as newDB
  component "New Module" as newModule
}

user -> stranglerModule: Use Strangler Features

stranglerModule --> legacyApp: Invoke legacy functionality
stranglerModule --> newModule: Delegate requests

stranglerModule --> stranglerDB: Read/Write Data
newModule --> newDB: Read/Write Data

@enduml
```


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

In many cases, it's good to keep the outputs and behaviors of both systems identical at the beginning — 
this makes it easier to compare results between the old and new versions to ensure that the new one works properly.

Try to get rid of the legacy part of the application ASAP and implement new features in new system after it's deployed.
Don't rewrite the entire application at once — find small parts that can be migrated, and do it step by step.

The alternative is a never-ending story of having "the new system" that will solve everything but is in development for years.
It's very, very common for this to happen: Avoid this at all costs because it never ends well.

</span>
	</div>
	</div>

The integration can be done via synchronous calls, but that may be risky when integrating with legacy systems.
If you don't need a synchronous response to the public API from your legacy system, it's a good iea to integrate over Pub/Sub.

The previous diagram is, of course, just a general idea.
It rarely looks exactly like that in real life — it's about illustrating the high-level concept.
You should adjust it to your needs.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

It's not a requirement that components from the diagram be separate services/microservices.
The only requirement is good modularization, but that can be achieved in a single monolith if it works best for you.

</span>
	</div>
	</div>

This is how it will look like in our case:


```plantuml
@startuml

skinparam monochrome true

!define ICONURL https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/v2.5.0/dist

!includeurl ICONURL/C4.puml

title Strangler Pattern

actor "User" as user

node "Dead Nation's\nPOST /book-tickets" as legacy {
  component "Legacy Application" as legacyApp
}

node "POST /book-tickets" as strangler {
  database "Database" as stranglerDB
  component "HTTP handler" as stranglerModule
}

node "POST /tickets-status" as newApp {
  database "Database" as newDB
  component "HTTP Handler" as newModule
}

user -> stranglerModule: Use Strangler Features

stranglerModule --> legacyApp: Invoke legacy functionality
legacyApp ->> newModule: Delegate requests
stranglerModule ..> newModule: Call directly (after removing Dead Nation)

stranglerModule --> stranglerDB: Read/Write Data
newModule --> newDB: Read/Write Data

@enduml
```

We would like to introduce a new endpoint: `POST /book-tickets`.
This is the endpoint that is called by our frontend when the user is trying to book a ticket.
Now it's handled by Dead Nation — we want to catch this request and let Dead Nation call `POST /tickets-status`.
After replacing Dead Nation, we will call `POST /tickets-status` directly.

The high-level idea of the Strangler Pattern is maintained: We have a part of our application that catches all requests and delegates them to the legacy system.
The difference, in our scenario, is that we won't remove this component after migration: It will just be a new endpoint.

That's the high-level idea. In the next few exercises, we'll do a step-by-step migration.
