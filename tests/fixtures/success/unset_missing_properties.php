<?php
unset($missing->x); $n = null; unset($n->x, $n[0]); $o = new stdClass; unset($o->x);
echo isset($missing) + 0, empty($n) + 0, empty($o->x) + 0, "\n";
