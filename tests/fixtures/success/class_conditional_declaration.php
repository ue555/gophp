<?php
if (true) { class Active { function value() { return 8; } } }
if (false) { class Unselected { function value() { echo 'bad'; return 1 / 0; } } }
$o = new Active; echo $o->value(), "\n";
