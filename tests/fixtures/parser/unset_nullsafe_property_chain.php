<?php
echo "before\n";
if (false) { unset($x?->item->field); }
echo "after\n";
