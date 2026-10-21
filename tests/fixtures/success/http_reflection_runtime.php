<?php
trait HttpNamedTrait {}
class HttpNamedClass {
    use HttpNamedTrait;
    public function __construct(public int $value = 4) {}
}
$reflection = new ReflectionClass(HttpNamedClass::class);
$instance = $reflection->newInstanceArgs([7]);
echo $instance->value, '|', count($reflection->getTraits()), '|', $reflection->isCloneable(), '|', $reflection->getConstructor()->isConstructor();
