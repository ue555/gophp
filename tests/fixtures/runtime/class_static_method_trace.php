<?php
class StaticTrace { public static function fail() { return 1 / 0; } }
StaticTrace::fail();
