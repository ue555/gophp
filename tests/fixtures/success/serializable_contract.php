<?php
$reflection = new ReflectionClass('Serializable');
echo $reflection->hasMethod('serialize') ? 'serialize' : 'missing', "\n";
echo $reflection->hasMethod('unserialize') ? 'unserialize' : 'missing', "\n";

abstract class Deferred implements Serializable {}
class Complete extends Deferred {
    public function serialize() { return 'payload'; }
    public function unserialize(string $data) { echo $data, "\n"; }
    public function __serialize(): array { return []; }
    public function __unserialize(array $data): void {}
}
$object = new Complete();
echo $object->serialize(), "\n";
$object->unserialize('restored');
