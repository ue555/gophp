<?php
// class_implements() order for user-defined classes: the parent class's list
// first, then the interfaces declared on the class in order, then each declared
// interface's own inherited interfaces, skipping names already listed.
interface Z {}
interface B extends Z {}
interface A extends B {}
interface D {}
interface F extends D, Z {}

class C implements A, D {}
class E extends C implements Countable
{
    public function count(): int
    {
        return 0;
    }
}
class G extends E implements F, Stringable
{
    public function __toString(): string
    {
        return 'g';
    }
}

foreach (['A', 'F', 'C', 'E', 'G'] as $name) {
    echo $name, ': ', implode(',', class_implements($name)), "\n";
}
echo implode(',', array_keys(class_implements(new G()))), "\n";
echo json_encode(class_implements('Z')), '|', json_encode(class_implements(new ArrayIterator([])) !== []), "\n";
