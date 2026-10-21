<?php
namespace NewNames;

class Thing {
    public function value() { return 3; }
}

function make() { return new \NewNames\Thing; }

$object = new \stdClass;
$object->value = 4;
$another = new
\NewNames\Thing;
$commented = new /* keyword boundary */ \NewNames\Thing;

try {
    throw new \Exception('boom');
} catch (\Exception $e) {
    echo make()->value(), '|', $object->value, '|', $another->value(), '|', $commented->value(), '|', $e->getMessage(), "\n";
}
