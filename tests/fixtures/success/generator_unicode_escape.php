<?php
echo "A\u{00A0}B" === "A\xc2\xa0B" ? 'yes|' : 'no|';
$name = 'N';
echo "{$name}\u{1F600}";
