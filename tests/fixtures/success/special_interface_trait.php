<?php
interface SpecialContract { public function identity(self $value): self; }
trait SpecialIdentity { public function identity(SpecialContract $value): self { return $this; } }
class SpecialImplementation implements SpecialContract { use SpecialIdentity; }
$value = new SpecialImplementation;
echo $value->identity($value) === $value ? 'same' : 'wrong';
