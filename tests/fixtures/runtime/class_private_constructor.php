<?php
class PrivateConstructor { private function __construct() { echo 'bad'; } }
$o = new PrivateConstructor($missing);
