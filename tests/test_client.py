# Mock the async send_request method
mock_tools = [{"name": "tool1"}, {"name": "tool2"}]
async def mock_send_request(*args, **kwargs):
    return {"result": {"tools": mock_tools}}
self.mock_session.send_request = AsyncMock(side_effect=mock_send_request) 