<?php
$first = [1, 'name' => 'old'];
$second = [...$first, 2, 'name' => 'new'];
echo $second[0], '|', $second[1], '|', $second['name'];
