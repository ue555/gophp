<?php
class BaseCounter {
    protected $value = 2;
    public function read() { return $this->value; }
}
class DoubleCounter extends BaseCounter {
    public function doubled() { return $this->value * 2; }
}
$counter = new DoubleCounter;
echo $counter->read(), '|', $counter->doubled(), "\n";
