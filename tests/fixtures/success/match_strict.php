<?php
$x = "1"; echo match ($x) { 1 => "int", "1" => "string", default => "other" }, "|";
echo match (2) { 1, 2 => "yes", default => 1 / 0 }, "\n";
