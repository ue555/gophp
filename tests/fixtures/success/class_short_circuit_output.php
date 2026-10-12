<?php
class Effects {
    public $n = 0;
    function run() { echo 'called'; $this->n = $this->n + 1; return true; }
    function fail() { echo 'bad'; return $missing / 0; }
}
$o = new Effects;
false && $o->fail(); true || $o->fail();
if (true) echo 'ok-'; elseif ($o->fail()) echo 'bad';
if (false) $o->run();
echo (true && $o->run()) + 0, ':', $o->n, "\n";
