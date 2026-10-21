<?php
$doc = new DOMDocument;
$doc->loadHTML('<html><body><div class="notice"><span>A</span><span>B</span></div></body></html>', LIBXML_NOERROR | LIBXML_NOBLANKS);
$body = $doc->getElementsByTagName('body')->item(0);
$spans = $doc->getElementsByTagName('span');
echo $body->nodeName, '|', $spans->length, '|';
foreach ($spans as $span) {
    echo $span->nodeValue;
}
echo '|', $body->childNodes[0]->getAttribute('class'), '|';
echo $spans->item(0)->nextSibling->nodeValue, '|', $doc->saveXML($spans->item(0));
