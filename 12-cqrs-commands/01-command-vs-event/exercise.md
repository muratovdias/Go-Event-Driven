# Commands vs Events

In [the Events module](/trainings/go-event-driven/exercise/9dd9d5ea-00e5-45ca-a4ca-0bd53ba42d7d), we said that **you should not use passive-aggressive
events.**
You should use events to model things that happened and represent them as immutable facts and not care
about what the subscriber does with the events.

However, what if you want to model something that *should* happen, and you care about a specific action to be performed?
**For example, suppose you want to send a notification to a user.**
You want to be sure that this operation will be eventually executed, but you don't want to wait for the result synchronously.
Commands are perfect for such scenarios.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Like events, commands are not intended to return any result to the emitter and are meant to be processed asynchronously.
What if you expect a specific action to be performed, but you want to know the result and do it synchronously?
This is a good case for a good old RPC or HTTP call.

</span>
	</div>
	</div>

## Commands

Commands are perfect for scenarios where you want to do something asynchronously.
While with events you are emitting the events and don't care what happens with them, in the case of commands, you
expect some reaction.

Usually, unlike events, commands are consumed by only one consumer.

While it's mostly a naming thing, in bigger systems, it's helping in understanding how the system works and what 
the expected behaviors are.
Usually, the commands' message broker may be a bit different — you don't want to have many consumers for one command.
You may also want to use different monitoring strategies for commands and events.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

You may have heard about Command Query Responsibility Segregation (CQRS).
One of the most common misunderstandings about CQRS is that commands need to be handled asynchronously.
That's not true!

For the sake of the training, we are using commands as asynchronous messages, but that's not required.
The only requirement of CQRS is that you should separate write models from read models.

Synchronous commands are out of scope of this training. 
You can read about that approach in our [Introducing basic CQRS by refactoring a Go project](https://threedots.tech/post/basic-cqrs-in-go/) article.

</span>
	</div>
	</div>

Watermill provides `CommandBus` and the `CommandProcessor` interface, which are analogous to `EventBus` and `EventProcessor`.


## Exercise

File: `12-cqrs-commands/01-command-vs-event/main.go`

Replace the `NotificationShouldBeSent` event with the `SendNotification` command.
You should also replace `cqrs.EventProcessor` with `cqrs.CommandProcessor`.

Commands should be published to the `commands` topic.
