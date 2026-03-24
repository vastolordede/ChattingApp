import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      home: const PingPage(),
    );
  }
}

class PingPage extends StatefulWidget {
  const PingPage({super.key});

  @override
  State<PingPage> createState() => _PingPageState();
}

class _PingPageState extends State<PingPage> {
  String result = 'Chưa gọi API';

  Future<void> callPingApi() async {
    try {
      final response = await http.get(
        Uri.parse('http://10.0.2.2:8080/ping'),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        setState(() {
          result = data['message'] ?? 'Không có message';
        });
      } else {
        setState(() {
          result = 'Lỗi HTTP: ${response.statusCode}';
        });
      }
    } catch (e) {
      setState(() {
        result = 'Không kết nối được BE: $e';
      });
    }
  }

  @override
  void initState() {
    super.initState();
    callPingApi();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Test FE ↔ BE'),
      ),
      body: Center(
        child: Text(
          result,
          style: const TextStyle(fontSize: 24),
        ),
      ),
    );
  }
}