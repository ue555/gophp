<?php
interface NeedTwoAgain { public function load($first, $second); }
interface NeedOneAgain { public function load($first); }
interface NeedsBothAgain extends NeedTwoAgain, NeedOneAgain {}
class Incomplete implements NeedsBothAgain {
    public function load($first) {}
}
