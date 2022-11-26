# Problem 3

This is the first non-trivial problem.
The program must handle the method "evaluate".
The params contain:

* `expression`: A string which is a valid [Lambda calculus](https://en.wikipedia.org/wiki/Lambda_calculus) expression.

The response must contain:

* `expression`: The result of evaluating the given expression in [Applicative Order](https://en.wikipedia.org/wiki/Evaluation_strategy#Strict_evaluation).

The syntax and evaluation method will be described in details below.

If the answer is correct, you will receive a link to the next problem.
Note that besides examples in this document, the grading script will also serve as a collection of test cases.

# Syntax
## Overview

The string must be a valid `Expression`, which can take one of the three forms:

* `Variable`
* `Abstraction`
* `Application`

## Variable

A `Variable` is a single lowercase alphabet character.
For example, these are valid `Variable` expressions: `a`, `b`, `c`...
These are invalid `Variable`:

* `hello`: More than one character
* `(`: Not an alphabet

## Abstraction

An `Abstraction` has the following form: `!<Variable>.<Expression>`.
This reads:

* Starts with `!` (Exclamation mark, ascii: 33)
* Followed by a `Variable`, defined in the previous section.
* Followed by: `.` (Period, ascii: 46)
* Followed by an `Expression`.

The `Variable` after the `!` character is also referred to as the argument.
The `Expression` after the `.` character is also referred to as the body.

Examples of an `Abstraction` expression:

* `!x.x`
* `!x.y`
* `!x.!y.y`: Because `!y.y` is an expression of type `Abstraction`, `!x.!y.y` is also an expression.
  Take note that the body can be any arbitrary expression, not just a `Variable`.
* `!x.(a b)`: The body is an `Application` defined in the next section.

These are invalid:

* `!(a b).c`: The argument must be a `Variable`, which is defined a single lowercase alphabet character.
* `!b. d`: Random empty space
* `!b.c `: Trailing space

## Application

An `Application` has the following form: `(<Expression> <Expression>)`.
This reads:

* Starts with `(` (Left parenthesis, ascii: 40)
* Followed by any `Expression`
* Followed by a white space: ` ` (ascii: 32)
* Followed by another `Expression`
* Followed by `)` (Right parenthesis, ascii: 41)

The first `Expression` is also referred to as the left hand side (LHS).
The second `Expression` is also referred to as the right hand side (RHS).

Examples of an `Application` expression:

* `(a b)`
* `((!z.z b) !x.(b c))`: Application of sub-expressions.
  Notice that the definition allows both LHS and RHS to be any arbitrary expressions.

These are invalid:

* `a  b`: No enclosing parentheses
* `(a b) `: Trailing space
* `(a b`: No right parenthesis

## Summary

There are 3 forms for an `Expression`.
Complex forms such as `Abstraction` and `Application` can contain nested expressions due to their recursive definition.
All forms can still be collectively referred to as "expression".

The syntax for this problem is much stricter than often seen elsewhere.
This should make it a lot easier to parse.
Remember that you are not expected to handle invalid input.

# Evaluating an expression

It is basically evaluation in [applicative order](https://en.wikipedia.org/wiki/Evaluation_strategy#Strict_evaluation).
Familiar readers can immediately start solving this problem without reading further.

We will define a procedure called `evaluate` which accepts an `Expression` and return the result which is also an `Expression`.
The signature of this procedure can be written as: `evaluate(expression: Expression): Expression`.
Since the expression can take one of 3 forms, the procedure is also split into 3 cases.

## Evaluating Variable

A `Variable` evaluates into itself.
For example: `x` evaluates into `x`.

## Evaluating Abstraction

An `Abstraction` evaluates into itself.
For example: `!x.y` evaluates into `!x.y`.

## Evaluating Application

`Application` is where things get complicated.
To `evaluate` an `Application`, do the following steps:

1. `evaluate` the LHS, call it `left`
2. `evaluate` the RHS, call it `right`
3. There are 2 cases:

    1. If `left` is **NOT** an `Abstraction`, return `(<left> <right>)`.
      This is an `Application` where the `LHS` and the `RHS` are the results from the corresponding `evaluation`.
    2. If `left` is an `Abstraction` with body `body` and argument `arg`, substitute all occurences of `arg` in `body` with `right`.
      Then `evaluate` this result again.

To understand substitution, one must first understand bound and free variable.

### Bound (and free) variables

A `Variable` is either "bound" or "free".
Intuitively, one could think of a bound variable as references to a function argument.
A free variable is either undefined or defined elsewhere as a global variable.

For example, given the following piece of Javascript:

```javascript
x( // This `x` is free
	(x) => { // `x` as the argument of this lambda function
		return x; // This `x` is bound to the `x` in the argument
	}
);
```

In fact, the above piece of code would be equivalent to: `(x !x.x)`.
Even though the LHS of this expression also contains `x`, this `x` is **NOT** the same as the `x` in the RHS.
One is free, the other is bound.

The formal definition is as follow:

* If the expression is a `Variable`, the variable is free in this expression.
* If the expression is an `Application`, the union of the free variables in the LHS and the RHS are the free variables in this expression.
* If the expression is an `Abstraction`, all the free variables in the body that do not have the same name as the argument are free in this expression.

Going back to the previous example, `(x !x.x)`.
There is no free variable in `!x.x` (RHS) but `x` is free in `x` (LHS).
Therefore, `x` (specifically, the LHS `x`) is free in the application `(x !x.x)`.
This can be shown in the following diagram:

```
(x !x.x)
    ^ â”‚
    â””â”€â”˜
```

The arrow goes from a variable to the argument that it is bound to.

For another example, consider this expression: `!y.!x.(y x)`.

```
!y.!x.(y x)
 ^  ^  â”‚ â”‚
 â”‚  â”‚  â”‚ â”‚ 
 â”‚  â””â”€â”€â”¼â”€â”˜
 â””â”€â”€â”€â”€â”€â”˜ 
```

In `!x.(y x)`, `x` is bound (to the argument) and `y` is free.
However, in `!y.!x.(y x)`, `y` is bound to the outer `y` due to having the same name.
Therefore, there is no free variable in the abstraction `!y.!x.(y x)`.

Binding does not have to be one-to-one.
In `!x.(x x)`, both `x`-s in the body are bound to the argument:

```
!x.(x x)
 ^  â”‚ â”‚
 â”‚  â”‚ â”‚  
 â”œâ”€â”€â”˜ â”‚
 â””â”€â”€â”€â”€â”˜ 
```

Intuitively, just think about the scoping (and shadowing) rule in most programming languages.

### Substitution

For brevity, we use `{v => e1}[e2]` to denote: "Substitute all occurences of variable `v` in expression `e2` with expression `e1`".

This forms the basis for evaluating an application.
Intuitively, "application" can be understood as function call in most programming languages.

Consider this piece of Javascript:

```javascript
((x) => x + 1)(1)
```

We can substitute `x` in the body of `(x) => x + 1` with `1` to get `1 + 1`.

However, substitution cannot be done naively.
Recall that in the previous section, we have learnt that variables can be bound or free, consider: `(!x.(x !x.x) y)`.
A naive substitution of all `x` in `(x !x.x)` with `y` would result in: `(y !x.y)` instead of `(y !x.x)`.
Notice how the binding of variables have changed.

Substitution is also split in to 3 cases.

#### Substituting into Variable

`{v1 => e1}[v2]`:

* If `v1 == v2` then `{v1 => e1}[v2] = e1`.
* If `v1 != v2` then `{v1 => e1}[v2] = v2`.

Intuitively, you can't replace something that is not there.

Examples:

* `{x => z}[x] = z`: First rule, `v1 == v2`
* `{x => z}[y] = y`: Second rule, `v1 != v2`

#### Substituting into Application

`{v => e}[(lhs rhs)] = ({v => e}[lhs] {v => e}[rhs])`.

Intuitively, distribute the substitution to the sub expressions.

Examples:

```
{x => z}[(x y)]
= ({x => z}[x] {x => z}[y])
= (z y)
```

#### Substituting into Abstraction

`{v1 => e1}[!v2.e2]`

There are 3 cases:

* If `v1 == v2` then `{v1 => e1}[!v2.e2] = !v2.e2`.

  Intuitively, the variable is already bound and can't be substituted.
* If `v1 != v2` and `v2` is **NOT** a free variable name in `e1` then `{v1 => e1}[!v2.e2] = !v2.{v1 => e1}[e2]`.

  Intuitively, the substitution can be pushed down without risk of name conflicts.
* If `v1 != v2` and `v2` is a free variable name in `e1` then we need to create a fresh variable name `v3` that is not yet used anywhere.
  `{v1 => e1}[!v2.e2] = {v1 => e1}[!v3.{v2 => v3}[e2]]`.

  Intuitively, to avoid a free variable in `e1` from being bound after substitution, we must rename the argument of the abstraction before proceeding.
  This does not change the nature of the abstraction.
  For example: `(x) => console.log(x)` is equivalent to `(y) => console.log(y)` even when we rename `x` to `y` becacause all bound instances are also renamed.

Examples:

* `{x => y}[!x.x] = !x.x`: First rule, `v1 == v2`
* `{x => y}[!z.x] = !z.y`: Second rule, `v1 != v2` and `v2` (`z`) is **NOT** free in `e1` (`y`)

An example of applying the third rule:

```
{x => y}[!y.(x y)]
= {x => y}[!z.{y => z}[(x y)]]
= {x => y}[!z.(x z)]
= !z.{x => y}[(x z)]
= !z.(y z)
```

#### Evaluation of Abstraction revisited

Step `3` in [Evaluating Application](#evaluating-application) should now read:

> There are 2 cases:
>
> 1. If `left` is **NOT** an `Abstraction`, return `(<left> <right>)`.
> 2. If `left` is an `Abstraction` with body `body` and argument `arg`, create `exp = {arg => right}[body]`.
>   Then `evaluate(exp)`.

#### Note on fresh variables

Only these known requirements are enforced:

* Fresh variables must not be used before
* Variable names must be a single lowercase alphabet character as defined in the syntax

You are free to choose any character to be used as fresh variables during substitution.
The grading script understands expression equivalence.
For example, it treats `!x.x` to be the same expression as `!y.y`.