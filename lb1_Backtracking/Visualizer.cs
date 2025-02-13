using System;
using System.Collections.Generic;
using SkiaSharp;
using System.Linq;

namespace lb1_Backtracking;

public class Visualizer
{
    public void VisualizePartition(List<int[]> squarePositions, int gridSize, string filename)
    {
        const int resolution = 800;
        int squareSize = (int)Math.Round(resolution / (double)gridSize);
        var rand = new Random();
        
        using (var bitmap = new SKBitmap(resolution, resolution))
        using (var canvas = new SKCanvas(bitmap))
        {
            canvas.Clear(SKColors.White);
            
            foreach (var square in squarePositions)
            {
                int x = square[0] * squareSize;
                int y = square[1] * squareSize;
                int size = square[2] * squareSize;
                
                using (var paint = new SKPaint { Style = SKPaintStyle.Fill, Color = new SKColor((byte)rand.Next(256), (byte)rand.Next(256), (byte)rand.Next(256)) })
                {
                    canvas.DrawRect(new SKRect(x, y, x + size, y + size), paint);
                }
                
                using (var borderPaint = new SKPaint { Style = SKPaintStyle.Stroke, Color = SKColors.Black, StrokeWidth = 3 })
                {
                    canvas.DrawRect(new SKRect(x, y, x + size, y + size), borderPaint);
                }
            }
            SaveBitmapToPng(bitmap, filename);
        }
    }

    private void SaveBitmapToPng(SKBitmap bitmap, string filename)
    {
        var filePath = Path.Combine(Environment.CurrentDirectory, filename);
        using (var image = SKImage.FromBitmap(bitmap))
        using (var data = image.Encode(SKEncodedImageFormat.Png, 100))
        using (var stream = File.OpenWrite(filePath))
        {
            data.SaveTo(stream);
        }
    }
}