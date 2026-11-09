<?php
interface A {}
interface B extends A {}
class P implements A {}
class C extends P implements B {}

$checks = [
    ['C', 'A'], ['C', 'B'], ['C', 'P'], ['C', 'C'],
    ['P', 'A'], ['P', 'P'], ['B', 'A'], ['B', 'B'], ['A', 'A'],
];
foreach ($checks as $check) {
    $class = new ReflectionClass($check[0]);
    echo $check[0], '->', $check[1], '=', $class->isSubclassOf($check[1]) ? 'true' : 'false', "\n";
}

$builtins = [
    ['ArrayIterator', 'Traversable'],
    ['ArrayIterator', 'Iterator'],
    ['ArrayIterator', 'Countable'],
    ['ArrayObject', 'IteratorAggregate'],
    ['ReflectionClass', 'Reflector'],
    ['ReflectionClass', 'Stringable'],
];
foreach ($builtins as $check) {
    $class = new ReflectionClass($check[0]);
    echo $check[0], '->', $check[1], '=', $class->isSubclassOf($check[1]) ? 'true' : 'false', "\n";
}
