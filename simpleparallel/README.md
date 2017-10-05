# Sample usage:

```
$ simpleparallel -cmd 'echo lol' -count 10 -jobs 6
lol

lol

lol

lol

lol

lol

lol

lol

lol

lol
```

or

```
$ seq 10 | simpleparallel -cmd 'echo {}' -jobs 6
1

6

2

7

8

10

4

3

9

5
```
