# Node To-Dos

## Function / Code

*   **What it does:** Runs an arbitrary snippet of logic on each record.
*   **Use case:** Any time you need custom, ad-hoc data munging or calculations that don’t fit a standard node.
*   **Go equivalent:** You already have `CodeNode`, which takes a `func([]map[string]interface{}) []map[string]interface{}`.

## Set / Resolve

*   **What it does:** Adds, updates, or removes fields on each item.
*   **Use case:** Normalizing field names, filling in default values, flattening nested JSON, or extracting sub-fields.
*   **Go implementation hint:** A small node that loops over input, modifies the maps in place (or returns new maps).

## Merge

*   **What it does:** Joins two (or more) input streams based on a key or index.
*   **Use case:** When you fetch extra data from two different APIs and need to combine them by a common ID.
*   **Go tip:** Buffer items from one stream in a map, then when items from the other stream arrive, look up and merge.

## SplitInBatches

*   **What it does:** Breaks a large array of items into chunks of N items.
*   **Use case:** Avoiding API rate-limit errors by sending only e.g. 10 requests at a time.
*   **Go implementation:** Write a node that, given an input slice of 100, emits 10 batches of size 10.

## Webhook

*   **What it does:** Exposes an HTTP endpoint that can kickoff a workflow run with incoming payload.
*   **Use case:** Integrating with Stripe, GitHub, or any service that calls webhooks.
*   **Go equivalent:** Wrap your `Workflow.Run` behind an `http.HandleFunc`, reading `r.Body` into your `ManualTrigger`.

## Cron / Schedule Trigger

*   **What it does:** Automatically starts a workflow at fixed intervals or cron expressions.
*   **Use case:** Polling an API every hour, doing nightly reports.
*   **Go approach:** Use a library like `robfig/cron` or a simple `time.Ticker` in a goroutine to call your trigger node.

## Webhook Response

*   **What it does:** Sends an HTTP response back to the caller, optionally including the output of the workflow.
*   **Use case:** “Query this endpoint and get back the transformed data.”
*   **Go hint:** Either block until `wf.Run` returns and write outputs, or respond immediately (202 Accepted) and stream back logs later.

## Switch / If

*   **What it does:** Routes items into different branches based on conditions.
*   **Use case:** If status == "pending", send to one API; if status == "approved", send to another.
*   **Go pattern:** Inspect each item, then append it only to the matching child’s input slice.

## WaitFor / Delay Until

*   **What it does:** Pauses processing of an individual item until either a time has passed or some external flag flips.
*   **Use case:** Send a follow-up email exactly 2 days after the first outreach.
*   **Go idea:** Check `item["next_send_timestamp"]`, sleep until then (or persist state externally and re-enqueue).

## MergeByKey

*   **What it does:** Groups incoming items by a key, waits until a count or timer threshold, then outputs one merged record.
*   **Use case:** Batch together all events for a user within a 5-minute window.
*   **Go sketch:** Maintain a map from key → slice of records, with a separate ticker per key to flush when time’s up.

## Error Trigger

*   **What it does:** Catches failures in downstream nodes and routes those items into an “error path” for retries or notifications.
*   **Use case:** If the LinkedIn invite API fails, send an alert email instead of crashing the whole workflow.
*   **Go approach:** Wrap each `node.Execute` in a `defer/recover` or `if err != nil` that appends the item to a separate error connection.