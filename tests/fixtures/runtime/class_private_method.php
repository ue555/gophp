<?php
class PrivateMethod { private function hidden() { echo 'bad'; } }
$o = new PrivateMethod;
$o->hidden($missing);
