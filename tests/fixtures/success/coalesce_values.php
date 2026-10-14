<?php
$a = ["x" => null, "y" => 3];
echo $missing ?? "m", "|", $a["x"] ?? "n", "|", $a["y"] ?? 4, "|", null ?? 5 ?? 6, "\n";
