<?php
$data = ['data' => 10, 'later' => 20];
$saved =& $data;
echo extract($data, EXTR_REFS), '|', $data, '|', $later, '|', $saved['data'], '|';
$data = 11;
$later = 21;
echo $saved['data'], '|', $saved['later'];
