<?php
interface RequiredLater { public function run(); }
abstract class DefersLater implements RequiredLater {}
echo "before\n";
if (true) { class ConcreteLater extends DefersLater {} }
echo "after\n";
