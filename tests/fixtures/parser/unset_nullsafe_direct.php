<?php
echo "before\n";
if (false) { unset($x?->item); }
echo "after\n";
