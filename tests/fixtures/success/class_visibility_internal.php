<?php
class Secret {
    private $x = 2;
    protected $y = 3;
    private function hidden() { return $this->x + $this->y; }
    protected function other($o) { return $o->x; }
    public function run($o) { return $this->hidden() + $this->other($o); }
}
$a = new Secret; $b = new Secret;
echo $a->run($b), "\n";
