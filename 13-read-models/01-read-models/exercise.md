# Read Models


<div class="alert alert-dismissible bg-info text-white d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-white">
			Background	
		</h3>
        <span>

Like it usually happens, during migrating out from Dead Nation we forgot about one more thing -
**Operations team had very nice dashboard in the Dead Nation website.**
They were using it as central place for doing their job.
It gave them view on all bookings, tickets and receipts.

The exact list of information that they used was:

- When the booking was created,
- List of tickets for the booking, with:
    - Price,
    - Customer email,
    - When the ticket was confirmed,
    - When the ticket was refunded (if it was),
    - When the ticket was printed,
    - File name of printed ticket (so they can download it by hand),
    - When the receipt was issued,
    - The receipt number.

This is a blocker before getting rid of Dead Nation.
**We need to provide that data to operations via our API.**

</span>
	</div>
	</div>

If you take a look on the list of information requested by Ops Team, we are not storing all of them in the database.
One option would be adding those informations to already existing tables.
But there is also one alternative approach: the _read model_ pattern.

In this pattern, we are storing data in a format that is optimized for fast reading (for example in API or by some processes).
It should be a format that is closest to desired format and doesn't require any transformations
In our case, we will store in a database JSON with format same as format that we will return from the API.

**In non-complex systems it would be enough to just query all the data from a single database, joining all needed tables.**
In the project we work on now it would be a good approach.
But let's do some over-engineering to practice.

### When to use read models?

There are multiple scenarios, when read models are a good choice.

**Read models are a great choice places where you have a data with high read throughput.**
You can build a read model that is optimized for reading and scale your database horizontally.
It's not uncommon, to store read models in a different kind of database than your other data.
It's called [Polyglot persistence](https://martinfowler.com/bliki/PolyglotPersistence.html)


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

In our case, we could use Elasticsearch to store the read models (or another NoSQL database).
It could provide nice search and filtering capabilities for the Ops Team.

Running Elasticsearch would take most of the resources used for the training,
so we'll just stick to PostgreSQL.
But keep in mind that it's an option to consider in a real system.

</span>
	</div>
	</div>

**Read models also help to give teams more autonomy on how they are storing data.**
Sometimes, some teams have different requirements for format of stored data.
With read models, they can store data in a format that is optimized for their use case without affecting other teams.

**You can also use read models, when you want to migrate out of a legacy system and build a new model based on
events.**
It's not uncommon, that legacy systems depend on large data models, which are hard to maintain 
when they are divided to multiple, smaller services.
We can emit events from the legacy system and build a new read model based on them in decoupled way.

Another scenario, when Read Models can help you, is a situation when some read use cases (for example, API endpoints, but not only)
has very high resiliency requirements, and data is now spread across multiple data stores.
With Read Model, you can aggregate all the required data in one place.
It simplifies reading operations a lot, so it's easier to ensure the stability of this part of infrastructure.

Last, but not least: you can use Read Models as a cache.
Usually one of the biggest problems with building a cache, is the invalidation policy.
For how long it should be cached? When should it be invalidated?
You can use a read model as an alternative, which is updated by events and is always up-to-date.


### Read model vs write model and source of truth

To understand Read Models properly, it's important to understand the difference between _read model_ and _write model_.

**_Write model_ is the place, which we consider as the place that is non-eventually consistent and has
always the latest available data.
This is the place, which we are updating and which guarantees domain invariants.**
Your write model should be a _single source of truth_ in your application.
If you want to update the booking, you update the `bookings` table.
If you want to check how many tickets are available, you check the `tickets` table.
**Those tables are our _write models_ and _single source of truth_ in our project**.

A typical monolith non-event-driven application has just write models.
It's part of the reason why they are so hard to scale and maintain.
It's hard to solve all problems with a single model.

Where can we find our _write models_ in our application?
The update of booking is emitting events, which are consumed by multiple places in the system.
They can store the data from events in their databases in form of read models.
But if they want to do the update, they should do the update in `bookings` table, not on their read models.

**Having a _single source of truth_ simplifies your system logic a lot.**
You don't need to think, about the list of places where you need to update the data.
You have just one place where you need to do that, and the downstream consumers needs to adapt to the changes.

### Cost of using 

Each of our decisions has tradeoffs. It's not different with read models.

The first cost is the cost of extra used storage.
We will need to duplicate some data in multiple databases.
On the other hand, that's what we already do in other contexts. 
RAID, High-Availability, horizontal scaling - it's all about duplicating data. 

Another cost is making migrations harder.
Now, instead of updating one model, we need to update multiple of them.
It gives you one advantage though - you don't need to migrate all models at once.
Some of them may be exposed publicly, so you may not want to change them.
With a single model it's hard to do, while read models give you this choice.

The last complexity is that you need to update multiple places in case of data bugs.
Like with migrations, it's not enough to run an SQL update on a single database.
You need to emit events for all read models and update them (or run multiple SQL updates on multiple databases, but we don't recommend it).
