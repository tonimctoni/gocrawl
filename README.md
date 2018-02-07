Small web crawler written in go.

# todo

 - Quote: "Caller should close resp.Body when done reading from".

 - Make sure `strings` in StringReservoir has constant capacity.

 - Maybe don't use defer if there is only one return in a function.

 - Do not have the css stuff running in an extra dedicated thread.

 - Filter already used urls outside of `get_urls` function.

 - Rename PageContent.go.

 - Enforce max urls to be crawled per site and max urls sharing a host.

 - Explain in documentation that structs are threadsafe if they are.

 - Test `UrlFinder`.

 - `UrlFinder.get_urls` should return roots of urls found too.

 - Test `remove_contained_strings_from_slice`.