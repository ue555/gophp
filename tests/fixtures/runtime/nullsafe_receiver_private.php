<?php
class Box { private $secret = null; }
$o = new Box;
echo $o?->secret?->take() ?? "fallback";
