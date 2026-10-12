<?php
if (false) { function selected() { return 'bad'; } }
if (true) { function selected() { return 'ok'; } }
echo SELECTED(), ':';
function outer() { function nested($x) { return $x + 1; } }
outer(); echo nested(3), "\n";
