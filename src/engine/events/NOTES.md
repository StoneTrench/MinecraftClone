
# Introduction

Game engines are like an onion, they have **layers**, and it is best practice to keep outer layers only dependant on inner layer, and not vice-versa. However you still want communication to be able to flow from the dependencies to the dependents, without the inner layer having to know about the outer systems. This can be achieved with **events**[^1].

**Owned events** ("*events are properties on the objects that own them*") are better than a **global bus**, because it easier to debug it, since it makes subscribing to an event explicit? They make subscribers visible, with a global bus the dependencies are invisible[^1].

# Best practices

+ Events only fire on state changes.
+ If owned events are used, events should not fire on object construction (only applies to OOP?) because the subscribers may not be ready yet.


[^1]: https://cleangamearchitecture.com/event-systems/