# gt 

A tiny but powerful Go internationalization (i18n) library.

## Installation

```sh
$ go get github.com/melvinmt/gt
```

## Usage
```go
package gt

import (
    "fmt"
    "github.com/melvinmt/gt"
)

var g = &gt.Build{
    Index: gt.Strings{
        "homepage-title": {
            "en":    "Hello World!",
            "es":    "¡Hola mundo!",
            "zh-CN": "你好世界!",
        },
        "homepage-welcome": {
            "en":    "Welcome to %s, %s.",
            "tr-TR": "%s, %s'ya hoşgeldiniz.",
            "nl":    "Welkom bij %s, %s",
        },
        "dashboard-notice": {
            "en":    "Hello %s#name, you have a new message from %s#friend.",
            "tr-TR": "%s#name merhaba, %sfriend'den yeni bir mesaj var.",
        },
        "invoice": {
            "en":    "You need to pay %10.2f#amount dollars in %d#days days.",
            "pt-BR": "Você precisa pagar %10.2f#amount dólares em %d#days dias.",
        },
    },
    Origin: "en", // the language you'll be translating from
}

func main() {
    // Key based:
    g.SetTarget("es")
    s1 := g.T("homepage-title")
    fmt.Println(s1) // outputs: ¡Hola mundo!

    // String based:
    g.SetTarget("zh") // notice that it's not necessary to include the region
    s2 := g.T("Hello World!")
    fmt.Println(s2) // outputs: 你好世界!

    // Parse arguments:
    g.SetTarget("nl")
    s3 := g.T("Welcome to %s, %s.", "Github", "Melvin")
    fmt.Println(s3) // outputs: Welkom bij Github, Melvin

    // As you can see in the previous example, you can use regular printf verbs
    // to parse arguments. The problem that you're facing here is that the order
    // of words is different in some languages:

    g.SetOrigin("es-LA")
    g.SetTarget("tr-TR")
    fmt.Println(g.T("Bienvenidos a %s, %s.", "Github", "Melvin"))
    // This outputs: Github, Melvin'ya hoşgeldiniz. This is roughly translated as:
    // Welcome to Melvin, Github.  Which is NOT what you want. You can solve this with
    // tag notation.

    // Tag notation:
    g.SetOrigin("en")
    s4 := g.T("Hello %s#name, you have a new message from %s#friend.")
    fmt.Println(s4, "Melvin", "Alex")
    // Outputs: Melvin merhaba, Alex'den yeni bir mesaj var. 
    // Which is in a different order, but correctly translated.

    // You can use any legal printf verb in combination with tags:
    g.SetOrigin("pt-BR")
    g.SetTarget("en")
    s5 := g.T("Você precise pagar %10.2f#amount dólares em %d#days dias.")
    fmt.Println(s5, 499.95, 5) // outputs: You need to pay 499.95 dollars in 5 days.
}

```

## Error handling

T() always returns strings. You can use the more verbose Translate() method to handle errors.

```go
key := "dashboard-notice"
s, err := g.Translate(key, "Melvin", "Alex")
if err != nil {
  // do something
}
```

## Edge cases

It's not recommended to pass duplicate anonymous printf verbs to **gt**, e.g. `"%s, %s, %d"`. It will work when the target strings will keep the arguments in order, but when one language requires to format the string as `"%s %d %s"`, **gt** will fail because it doesn't know which `%s` to swap. You can easily solve this by tagging duplicate verbs: `"%s#tag1 %s#tag2 %d"`.

Even when **gt** fails, it will try to return the origin string with formatted arguments. In this way, even when a translation has failed (or simply doesn't exist yet), you can at least present something to the end user.

```go
g := &gt.Build{
    Index: gt.Strings{
        "incomplete": {
            "en": "Sorry %s, this string is not translated yet!",
        },
    },
    Origin: "en",
}
g.SetTarget("es")
s, err := g.Translate("incomplete", "John")
if err != nil {
    fmt.Println(s) // outputs: Sorry John, this string is not translated yet!
}
```

By default, Origin and Target are set to `"xx"` to prevent out of bound runtime errors.

## Feedback

I just started coding in Go a week ago (jan '13) and I'm still learning everyday. Please tell me when something's not solved in a idiomatic or optimal way and I'll change it (better yet, make a pull request)! This is not to say that this library isn't ready to be used, it passes all the tests in [gt_test.go](https://github.com/melvinmt/gt/blob/master/gt_test.go) and you should be able to use it in your projects.

## Docs

http://godoc.org/github.com/melvinmt/gt

## History

### *01/20/2013*: v1.0.0

- initial version of gt
- gt passes all tests 
- wrote the documentation
