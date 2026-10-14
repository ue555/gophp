<?php
class Maybe { public $x = 3; } $a = null; $b = new Maybe; echo isset($a?->x) ? 1 : 0, "|", isset($b?->x) ? 1 : 0, "\n";
