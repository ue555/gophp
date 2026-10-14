<?php
class Holder { public $x = null; } $h = new Holder; $count = 0;
function holder($h) { echo "h"; return $h; }
holder($h)->x ??= 6; echo "|", $h->x, "\n";
