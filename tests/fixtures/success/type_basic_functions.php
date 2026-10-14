<?php
function add(int $a, float $b): float { return $a + $b; }
function label(?string $s): string { return $s ?? "none"; }
echo add("2", 3), "|", label(null), "\n";
