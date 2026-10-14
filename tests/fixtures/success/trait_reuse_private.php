<?php
trait InnerValue { private $value = 1; private function read() { return $this->value; } }
trait ValueTrait { use InnerValue; public function get() { return $this->read(); } }
class One { use ValueTrait; public function set($v) { $this->value = $v; } }
class Two { use ValueTrait; }
$a = new One;
$b = new Two;
$a->set(9);
echo $a->get(), '|', $b->get(), "\n";
