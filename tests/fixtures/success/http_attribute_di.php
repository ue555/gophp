<?php
#[Attribute]
class InjectName {
    public function __construct(public string $value) {}
}
class NamedService {
    public function __construct(#[InjectName('answer')] public string $name) {}
}
$class = new ReflectionClass(NamedService::class);
$parameter = $class->getConstructor()->getParameters()[0];
$attribute = $parameter->getAttributes()[0];
$service = new NamedService($attribute->newInstance()->value);
echo $attribute->getName(), '|', $service->name;
