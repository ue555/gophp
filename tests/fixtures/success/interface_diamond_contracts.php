<?php
interface NeedTwo { public function load($first, $second = null); }
interface NeedOne { public function load($first); }
interface NeedsBoth extends NeedTwo, NeedOne {}
class Complete implements NeedsBoth {
    public function load($first, $second = null) { return $first . ':' . $second; }
}
echo (new Complete)->load('a', 'b'), "\n";
