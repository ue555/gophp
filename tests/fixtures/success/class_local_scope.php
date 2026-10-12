<?php
class Local {
    function run($x) { $x = $x + 1; $local = 7; return $x + $local; }
}
$x = 100; $local = 200; $o = new Local;
echo $o->run(1), ':', $o->run(2), ':', $x, ':', $local, "\n";
