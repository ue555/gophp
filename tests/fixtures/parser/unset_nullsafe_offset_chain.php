<?php
echo "before\n";
if (false) { unset($x?->items[0]); }
echo "after\n";
