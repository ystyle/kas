package tpl

const PAGE = `<!DOCTYPE html>
<html>
    <head>
        <meta content="KC2/1.160/cf63ffb/win" name="generator"/>
        <title>{{ .id }}</title>
    </head>
    <body>
        <div>
            <img style="width:800px;height:1132px;margin-left:0px;margin-top:74px;margin-right:0px;margin-bottom:74px;" src="{{ .image }}"/>
        </div>
    </body>
</html>`

const OPF = `<package version="2.0" xmlns="http://www.idpf.org/2007/opf" unique-identifier="{ {{- .id -}} }">
    <metadata xmlns:opf="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/">
        <meta content="comic" name="book-type"/>
        <meta content="true" name="zero-gutter"/>
        <meta content="true" name="zero-margin"/>
        <meta content="true" name="fixed-layout"/>
        <meta content="KindleComicCreator-1.0" name="generator"/>
        <dc:title>{{ .title }}</dc:title>
        <dc:language>zh</dc:language>
        <dc:creator>{{ .author }}</dc:creator>
        <dc:publisher>hcc(hcomic creator)</dc:publisher>
        <meta content="portrait" name="orientation-lock"/>
        <meta content="horizontal-lr" name="primary-writing-mode"/>
        <meta content="800x1280" name="original-resolution"/>
        <meta content="false" name="region-mag"/>
        <meta content="cover-image" name="cover"/>
        <dc:source>KC2/1.160/cf63ffb/win</dc:source>
    </metadata>
    <manifest>
        <item href="toc.ncx" id="ncx" media-type="application/x-dtbncx+xml"/>
        <item href="cover-image.jpg" id="cover-image" media-type="image/jpg"/>
{{- range .images}}
        <item href="{{ .filename }}" id="item-{{ .item }}" media-type="application/xhtml+xml"/>
{{- end}}
    </manifest>
    <spine toc="ncx">
{{- range .images}}
        <itemref idref="item-{{ .item }}" linear="yes"/>
{{- end}}
    </spine>
</package>`

const TOC = `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE ncx PUBLIC '-//NISO//DTD ncx 2005-1//EN' 'http://www.daisy.org/z3986/2005/ncx-2005-1.dtd'>
<ncx version="2005-1" xmlns="http://www.daisy.org/z3986/2005/ncx/" xml:lang="en-US">
    <head>
        <meta content="" name="dtb:uid"/>
        <meta content="" name="dtb:depth"/>
        <meta content="0" name="dtb:totalPageCount"/>
        <meta content="0" name="dtb:maxPageNumber"/>
        <meta content="true" name="generated"/>
    </head>
    <docTitle>
        <text/>
    </docTitle>
    <navMap>
{{- range .images}}
        <navPoint playOrder="{{ .i }}" id="toc-{{ .i }}">
            <navLabel>
                <text>Page-{{ .id }}</text>
            </navLabel>
            <content src="{{ .filename }}"/>
        </navPoint>
{{- end}}
    </navMap>
</ncx>`
