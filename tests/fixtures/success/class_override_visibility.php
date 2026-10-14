<?php
class OverrideBase {
    public $number = 1;
    protected $other = 2;
    public function other() { return $this->other; }
    private function secret() { return "base"; }
    public function callSecret() { return $this->secret(); }
    public function label() { return "base"; }
}
class OverrideChild extends OverrideBase {
    public $number = 3;
    public $other = 4;
    private function secret() { return "child"; }
    public function label() { return "child"; }
}
$child = new OverrideChild;
echo $child->number, "|", $child->other, "|", $child->other(), "|", $child->callSecret(), "|", $child->label(), "\n";
