<?php
class Holder { private $secret = 1; }
$o = new stdClass;
$h = new Holder;
echo isset($o?->missing) ? "yes" : "no";
echo "|", $o?->missing ?? "fallback";
echo "|", empty($o?->missing) ? "empty" : "full";
echo "|", isset($undefined?->missing) ? "yes" : "no";
echo "|", $undefined?->missing ?? "fallback";
echo "|", isset($h?->secret) ? "yes" : "no";
echo "|", $h?->secret ?? "fallback";
echo "|", empty($h?->secret) ? "empty" : "full", "\n";
