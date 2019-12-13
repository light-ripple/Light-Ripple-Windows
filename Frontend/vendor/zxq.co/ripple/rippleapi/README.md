# rippleapi

This is the source code for Ripple's API.

- Origin: https://git.zxq.co/ripple/rippleapi
- Mirror: https://github.com/osuripple/api

## Note to fellow developers: this is not how you do it!

The API is crammed with terrible design. First of all, it is not RESTful, and as you'll come to learn, designing an API in a RESTful manner is good because it helps to create consistent design (across your API and other APIs). It also quite simplifies many other things:

* In the API, to get a specific item, you need to do e.g. `/users?id=1009`. It's much more useful to have these in the URL path directly (`/users/1009`) for a number of reasons:
  * It simplifies checks (`/users/scores?id=1009` will require a check to see if an ID is present. `/users/:id/scores` doesn't really need a check, because `/users/scores` won't match)
  * It gives a "feel" of hierarchy
  * There is no multiple ways of addressing a specific user. There is a single way: IDs. In the Ripple API, you can specify an username instead of an ID to get a specific user, but this is prone to failure in the event of the user changing the username, whereas an ID cannot (should not) change.
* You can show error codes to the user using HTTP status codes. This way, you can present the resource to the user without any wrapper (such as an object giving an "ok" field or, like in the API, a `code` parameter), so the user can likely reuse other parts for error handling that they already use for other http requests.
* GET merely shows a resource, is cacheable and "idempotent". This helps debugging (repeating the same request twice will yield the same result, unless of course the data changes otherwise), caching (you can answer with Cache-Control headers, which browsers understand).

The not-making-it-RESTful was the biggest sin of the API. In itself, the API was a step into the right direction (it is MUCH better than the official osu! API), but nowhere close to how an API actually is, ideally. If you are building an API, I won't recommend you a book, but instead I will recommend you to see what [GitHub](https://developer.github.com/v3/) does, as they will have probably faced most problems that you have, and provided an answer already. If you're unsure, check other APIs: Discord, Slack, Twitter, Stripe, to name a few.

